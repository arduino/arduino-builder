void setup(){
  func()();
}

void loop(){}

void (*func())(){
  return setup;
}
