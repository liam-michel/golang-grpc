package service

import (
	"testing"
)

func TestSquareNumber(t *testing.T) {
	tests :=[]struct{
		name string
		input int32
		output int32
	}{
		{"Square of 1", 1, 1},
		{"Square of 2", 2, 4},
		{"Square of 3", 3, 9},
		{"Square of 4", 4, 16},
		{"Square of 5", 5, 25},
		{"Square of 6", 6, 36},
	}

	for _,tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			result, err := SquareNumber(tt.input)
			if err !=nil{
				t.Fatalf("Error: %v", err)

			}
			if result != tt.output{	
				t.Fatalf("Expected %v but got %v", tt.output, result)
			}
			if result == tt.output{	
				t.Logf("Expected %v and got %v", tt.output, result)
			}

		})
	}

}