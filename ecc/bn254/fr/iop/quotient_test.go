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

package iop

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// computes x₃ in h(x₁,x₂,x₃) = x₁^{2}*x₂ + x₃ - x₁^{3}
// from x₁ and x₂.
func computex3(x []fr.Element) fr.Element {

	var a, b fr.Element
	a.Square(&x[0]).Mul(&a, &x[1])
	b.Square(&x[0]).Mul(&b, &x[0])
	a.Sub(&b, &a)
	return a

}

func allocatePol(size int, form Form) *Polynomial {
	f := new(Polynomial)
	f.Coefficients = make([]fr.Element, size)
	f.Info = form
	return f
}

func TestQuotient(t *testing.T) {

	// create the multivariate polynomial h
	// h(x₁,x₂,x₃) = x₁^{2}*x₂ + x₃ - x₁^{3}
	nbInputs := 3
	h := make(multivariatePolynomial, nbInputs)

	h[0].coeff.SetOne()
	h[0].exponents = []int{2, 1, 0}

	h[1].coeff.SetOne()
	h[1].exponents = []int{0, 0, 1}

	h[2].coeff.SetOne().Neg(&h[2].coeff)
	h[2].exponents = []int{3, 0, 0}

	// create an instance (f_i) where h holds
	sizeSystem := 8
	form := Form{Basis: Lagrange, Status: Locked, Layout: Regular}
	f := make([]*Polynomial, 3)
	f[0] = allocatePol(sizeSystem, form)
	f[1] = allocatePol(sizeSystem, form)
	f[2] = allocatePol(sizeSystem, form)
	for i := 0; i < sizeSystem; i++ {

		f[0].Coefficients[i].SetRandom()
		f[1].Coefficients[i].SetRandom()
		tmp := computex3([]fr.Element{f[0].Coefficients[i], f[1].Coefficients[i]})
		f[2].Coefficients[i].Set(&tmp)

		x := []fr.Element{
			f[0].Coefficients[i],
			f[1].Coefficients[i],
			f[2].Coefficients[i],
		}
		tmp = h.evaluate(x)
		if !tmp.IsZero() {
			t.Fatal("system does not vanish on x^n-1")
		}
	}

	// compute the quotient q

	// checks that h(f_i) = (x^n-1)*q

}
