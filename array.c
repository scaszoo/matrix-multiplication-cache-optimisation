/**
gcc array.c -o test -mcmodel=medium && ./array 
for HU gruenau servers
 */

#include<stdio.h>
#include <time.h>
#include <stdlib.h>

#define x 1200
#define y 1200

// long to use up more memory
long matrix[y][x]; // column, row
long matrixTurned[y][x];
long matrixResult[y][x];
long matrixTurnedResult[y][x];


int main() {
	printf("Let's have some fun with matrix multiplication.\n");
	srand(time(NULL)); // Initialization, should only be called once.

	int i, j, k;
	long sum;

	// initialize both matrices
	for(i=0; i<y; i++) { // row
		for(j=0; j<x; j++) { // column
			long tmp = rand() % 100;
			matrix[i][j] = tmp;
			matrixTurned[j][i] = tmp;
		}
	}

	// calculate matrix X matrix
	printf("\n\ncalculate matrix X matrix\n");
	printf("Timestamp: %ld\n", time(NULL));
	for(i=0; i<y; i++) { // row
		for(j=0; j<x; j++) { // column
			sum = 0;
			for(k=0; k<x; k++) { // K
				sum += matrix[k][j] * matrix[i][k];
			}
			matrixResult[i][j] = sum;
		}
		
	}
	printf("Timestamp: %ld\n", time(NULL));

	// calculate matrix X matrixTurned
	printf("\n\ncalculate matrix X matrixTurned\n");
	printf("Timestamp: %ld\n", time(NULL));
	for(i=0; i<y; i++) { // row
		for(j=0; j<x; j++) { // column
			sum = 0;
			for(k=0; k<x; k++) { // K
				sum += matrix[i][k] * matrixTurned[j][k];
			}
			matrixTurnedResult[i][j] = sum;
		}
		
	}
	printf("Timestamp: %ld\n", time(NULL));

	// compare results
	for(i=0; i<y; i++) { // row
		for(j=0; j<x; j++) { // column
			if (matrixTurnedResult[i][j] != matrixResult[i][j]) {
				printf("CRASH BURN\n %ld %ld\n", matrixResult[i][j], matrixTurnedResult[i][j]);	
				return 0;
			}
		}
	}

	return 0;
}
