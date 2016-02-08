#include <Arduino.h> // Arduino 1.0

#include "other.h"

MyClass::MyClass() {
}

void MyClass::init ( Stream *stream ) {
    controllerStream = stream;
}
