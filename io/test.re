// Test code for lexical analyzer
var isPrime = fn(c) { // Returns true if c is prime.
    if (c % 2 == 0) {
        return false;
    }
    var d;
    d = 3;
    for (d != c) { // If c < 0 we're in trouble!
        if (c % d == 0) {
            return false;
        }
        d = d + 2;
    }
    var s = "This is a useless string literal";
    return true;
};

/* Find greatest 
common divisor. */
var gcd = fn(a, b) {
    var s = " This is an invalid string literal  /// error
    var m;
    var #23sd;                                    ///error
    m = b % a;
    if (m == 0) {
        return a;
    } else {
        return gcd(m, a);
    }
};

var fibonacci = fn(n) { // Find n-th Fibonacci number.
    if ((n == 1) || (n == 2)) {
        return 1;
    } else {
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
};

var average = fn(n1, n2, n3) {
    var sum = n1 + n2 + n3;
    return sum / 3;
};

// Additional features demonstration
var testArrays = fn() {
    var arr = [1, 2, 3, 4, 5];
    var firstElement = arr[0];
    arr[2] = 10;
    return arr;
};

var logicalOperators = fn() {
    var a = true;
    var b = false;
    var result1 = a && b;
    var result2 = a || b;
    var result3 = !a;
    return result2;
};

var multilineString = fn() {
    var message = `This is a multi-line
string that demonstrates
the backtick syntax in
relang language`;
    return message;
};

// Main function
var main = fn() {
    var num = 17;
    var isPrimeResult = isPrime(num);
    var fibResult = fibonacci(10);
    var gcdResult = gcd(56, 98);
    
    return "Lexical analysis test complete";
};

main();

/*
