#ifndef other__h
#define other__h

class MyClass {
  public:
    MyClass();
    void init ( Stream *controllerStream );

  private:
    Stream *controllerStream;
};
#endif