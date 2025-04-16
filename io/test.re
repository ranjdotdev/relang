// test.re

/*
 * Function to check if two values are equal.
 * It returns true if they are equal, otherwise returns false.
 */
 
var isEqual = fn(x, y) {
    if (x == y) {
        return true;
    } else {
        return false;
    }
};

// Variable declaration
var five = 5; var ten = 10;
var result = isEqual(five, ten);

// Print the result
print(result);