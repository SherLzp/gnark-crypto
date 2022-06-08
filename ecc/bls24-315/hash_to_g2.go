// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bls24315

import (
	"github.com/consensys/gnark-crypto/ecc/bls24-315/fp"
	"github.com/consensys/gnark-crypto/ecc/bls24-315/internal/fptower"
)

// mapToCurve2 implements the Shallue and van de Woestijne method, applicable to any elliptic curve in Weierstrass form
// No cofactor clearing or isogeny
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-hash-to-curve-14#appendix-F.1
func mapToCurve2(u *fptower.E4) G2Affine {
	var tv1, tv2, tv3, tv4 fptower.E4
	var x1, x2, x3, gx1, gx2, gx, x, y fptower.E4
	var one fptower.E4
	var gx1NotSquare, gx1SquareOrGx2Not int

	//constants
	//c1 = g(Z)
	//c2 = -Z / 2
	//c3 = sqrt(-g(Z) * (3 * Z² + 4 * A))     # sgn0(c3) MUST equal 0
	//c4 = -4 * g(Z) / (3 * Z² + 4 * A)

	//Z  = 2
	//c1 = 8 + 33596659215742140129637122214961435938020655891235557129669780043232055771106905062908445266866*x³
	//c2 = 39705142709513438335025689890408969744933502416914749335064285505637884093126342347073617133568
	//c3 = 34201884056882185778406405395330386747248400816830429692069788896408054379591483015871146287242 + 19460743984535740826215914161141042964572206362769336154303137946305037352047225221026480695590*x + 22542073731636879082798682200224885822407588647098269372934034390571474185157692939207072589882*x² + 39120870487753725729765241885585802642512708522196829619858744308140817355466861636413036946525*x³
	//c4 = 26470095139675625556683793260272646496622334944609832890042857003758589395417561564715744755710 + 2036161164590432735129522558482511268970948841893064068464835154135276107339812428055057288901*x³

	//TODO: Move outside function?
	Z := fptower.E4{
		B0: fp.Element{4181239655115521941, 6707528626421712891, 16500631436646597700, 7468167222847106173, 203155518471302199},
		B1: fp.Element{0},
		B2: fp.Element{0},
		B3: fp.Element{0},
	}
	c1 := fptower.E4{
		B0: fp.Element{597561764214734418, 17301118142370108904, 15453102953399245648, 17771897912064275607, 126821463998334011},
		B1: fp.Element{0},
		B2: fp.Element{0},
		B3: fp.Element{11675424027742205484, 705286907219923506, 17323892202765057035, 9478112032423100124, 163636478684360919},
	}
	c2 := fptower.E4{
		B0: fp.Element{11164601423358853174, 17475228851327880835, 18222098035255651149, 13126167188689647896, 69872393236067596},
		B1: fp.Element{0},
		B2: fp.Element{0},
		B3: fp.Element{0},
	}
	c3 := fptower.E4{
		B0: fp.Element{11487849815991738919, 14002903393647637345, 16849071880587728963, 14029988556784598762, 187510217357509286},
		B1: fp.Element{4987342894163522335, 1897321579306219091, 2455958775128937592, 12576725190445155719, 220100328698682968},
		B2: fp.Element{1499398348614592911, 18171467641400330674, 3384106693253235807, 9112576606476320265, 150392972236928314},
		B3: fp.Element{15906358451573163609, 14121029319544903442, 6350540516720830097, 6260920200926301033, 83197043233891257},
	}
	c4 := fptower.E4{
		B0: fp.Element{8637626912539497957, 1970041370999271347, 6348326826683034245, 5316201229387375453, 72026280315034460},
		B1: fp.Element{0},
		B2: fp.Element{0},
		B3: fp.Element{13781820325308083698, 2941236485365606384, 4926262172237512167, 7023134340203533526, 174054710400837955},
	}

	one.SetOne()

	tv1.Square(u)       //    1.  tv1 = u²
	tv1.Mul(&tv1, &c1)  //    2.  tv1 = tv1 * c1
	tv2.Add(&one, &tv1) //    3.  tv2 = 1 + tv1
	tv1.Sub(&one, &tv1) //    4.  tv1 = 1 - tv1
	tv3.Mul(&tv1, &tv2) //    5.  tv3 = tv1 * tv2

	tv3.Inverse(&tv3)   //    6.  tv3 = inv0(tv3)
	tv4.Mul(u, &tv1)    //    7.  tv4 = u * tv1
	tv4.Mul(&tv4, &tv3) //    8.  tv4 = tv4 * tv3
	tv4.Mul(&tv4, &c3)  //    9.  tv4 = tv4 * c3
	x1.Sub(&c2, &tv4)   //    10.  x1 = c2 - tv4

	gx1.Square(&x1) //    11. gx1 = x1²
	//TODO: Beware A ≠ 0
	//12. gx1 = gx1 + A
	gx1.Mul(&gx1, &x1)                 //    13. gx1 = gx1 * x1
	gx1.Add(&gx1, &bTwistCurveCoeff)   //    14. gx1 = gx1 + B
	gx1NotSquare = gx1.Legendre() >> 1 //    15.  e1 = is_square(gx1)
	// gx1NotSquare = 0 if gx1 is a square, -1 otherwise

	x2.Add(&c2, &tv4) //    16.  x2 = c2 + tv4
	gx2.Square(&x2)   //    17. gx2 = x2²
	//    18. gx2 = gx2 + A
	gx2.Mul(&gx2, &x2)               //    19. gx2 = gx2 * x2
	gx2.Add(&gx2, &bTwistCurveCoeff) //    20. gx2 = gx2 + B

	{
		gx2NotSquare := gx2.Legendre() >> 1              // gx2Square = 0 if gx2 is a square, -1 otherwise
		gx1SquareOrGx2Not = gx2NotSquare | ^gx1NotSquare //    21.  e2 = is_square(gx2) AND NOT e1   # Avoid short-circuit logic ops
	}

	x3.Square(&tv2)   //    22.  x3 = tv2²
	x3.Mul(&x3, &tv3) //    23.  x3 = x3 * tv3
	x3.Square(&x3)    //    24.  x3 = x3²
	x3.Mul(&x3, &c4)  //    25.  x3 = x3 * c4

	x3.Add(&x3, &Z)                  //    26.  x3 = x3 + Z
	x.Select(gx1NotSquare, &x1, &x3) //    27.   x = CMOV(x3, x1, e1)   # x = x1 if gx1 is square, else x = x3
	// Select x1 iff gx1 is square iff gx1NotSquare = 0
	x.Select(gx1SquareOrGx2Not, &x2, &x) //    28.   x = CMOV(x, x2, e2)    # x = x2 if gx2 is square and gx1 is not
	// Select x2 iff gx2 is square and gx1 is not, iff gx1SquareOrGx2Not = 0
	gx.Square(&x) //    29.  gx = x²
	//    30.  gx = gx + A

	gx.Mul(&gx, &x)                //    31.  gx = gx * x
	gx.Add(&gx, &bTwistCurveCoeff) //    32.  gx = gx + B

	y.Sqrt(&gx)                             //    33.   y = sqrt(gx)
	signsNotEqual := g2Sgn0(u) ^ g2Sgn0(&y) //    34.  e3 = sgn0(u) == sgn0(y)

	tv1.Neg(&y)
	y.Select(int(signsNotEqual), &y, &tv1) //    35.   y = CMOV(-y, y, e3)       # Select correct sign of y
	return G2Affine{x, y}
}

// g2Sgn0 is an algebraic substitute for the notion of sign in ordered fields
// Namely, every non-zero quadratic residue in a finite field of characteristic =/= 2 has exactly two square roots, one of each sign
// Taken from https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/ section 4.1
// The sign of an element is not obviously related to that of its Montgomery form
func g2Sgn0(z *fptower.E4) uint64 {

	nonMont := *z
	nonMont.FromMont()

	sign := uint64(0)
	zero := uint64(1)
	var signI uint64
	var zeroI uint64

	signI = nonMont.B0[0] % 2
	sign = sign | (zero & signI)

	zeroI = g1NotZero(&nonMont.B0)
	zeroI = 1 ^ (zeroI|-zeroI)>>63
	zero = zero & zeroI

	signI = nonMont.B1[0] % 2
	sign = sign | (zero & signI)

	zeroI = g1NotZero(&nonMont.B1)
	zeroI = 1 ^ (zeroI|-zeroI)>>63
	zero = zero & zeroI

	signI = nonMont.B2[0] % 2
	sign = sign | (zero & signI)

	zeroI = g1NotZero(&nonMont.B2)
	zeroI = 1 ^ (zeroI|-zeroI)>>63
	zero = zero & zeroI

	signI = nonMont.B3[0] % 2
	sign = sign | (zero & signI)

	return sign

}

// MapToG2 invokes the SVDW map, and guarantees that the result is in g2
func MapToG2(u fptower.E4) G2Affine {
	res := mapToCurve2(&u)
	res.ClearCofactor(&res)
	return res
}

// EncodeToG2 hashes a message to a point on the G2 curve using the SVDW map.
// It is faster than HashToG2, but the result is not uniformly distributed. Unsuitable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
//https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/#section-6.6.3
func EncodeToG2(msg, dst []byte) (G2Affine, error) {

	var res G2Affine
	u, err := hashToFp(msg, dst, 4)
	if err != nil {
		return res, err
	}

	res = mapToCurve2(&fptower.E4{
		0: u[0],
		1: u[1],
		2: u[2],
		3: u[3],
	})

	res.ClearCofactor(&res)
	return res, nil
}

// HashToG2 hashes a message to a point on the G2 curve using the SVDW map.
// Slower than EncodeToG2, but usable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#section-3
func HashToG2(msg, dst []byte) (G2Affine, error) {
	u, err := hashToFp(msg, dst, 2*4)
	if err != nil {
		return G2Affine{}, err
	}

	Q0 := mapToCurve2(&fptower.E4{
		0: u[0],
		1: u[1],
		2: u[2],
		3: u[3],
	})
	Q1 := mapToCurve2(&fptower.E4{
		0: u[4+0],
		1: u[4+1],
		2: u[4+2],
		3: u[4+3],
	})

	var _Q0, _Q1 G2Jac
	_Q0.FromAffine(&Q0)
	_Q1.FromAffine(&Q1).AddAssign(&_Q0)

	_Q1.ClearCofactor(&_Q1)

	Q1.FromJacobian(&_Q1)
	return Q1, nil
}

func g2NotZero(x *fptower.E4) uint64 {
	//Assuming G1 is over Fp and that if hashing is available for G2, it also is for G1
	return g1NotZero(&x.A0) | g1NotZero(&x.B1) | g1NotZero(&x.B2) | g1NotZero(&x.B3)

}