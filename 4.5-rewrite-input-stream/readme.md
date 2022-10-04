

# Rewriting the Input Stream

Let’s build a tool that processes Java source code to insert serialization identifiers, serialVersionUID.

Input:

```java
import java.util.List; 
import java.util.Map; 

public class Demo {
    void f(int x, String y) { }
    int[ ] g(/*no args*/) { return null; } 
    List<Map<String, Integer>>[] h() { return null; }
}
```


Output:

```java
import java.util.List; 
import java.util.Map; 

public class Demo {
	public static final long serialVersionUID = 1L;
    void f(int x, String y) { }
    int[ ] g(/*no args*/) { return null; } 
    List<Map<String, Integer>>[] h() { return null; }
}
```
