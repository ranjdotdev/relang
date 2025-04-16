var isEqual = fn(x, y) {
    if (x == y) {
        return true;
    } else {
        return false;
    }
};

var five = 5; var ten = 10;
var result = isEqual(five, ten);

print(result);