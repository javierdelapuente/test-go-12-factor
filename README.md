
Simple project written in Go, with the next endpoints:


 - /: Return "hello world"
 - /env: List all env variables (json).
 - /redis/status:
 - /mongodb/status:
 - /postgresql/status
 - /s3/status
 - /mysql/status
 - /config/<config_name>
 - /sleep
  

```mermaid
  graph TD;
      A-->B;
      A-->C;
      B-->D;
      C-->D;
```


```mermaid
classDiagram
    note "From Duck till Zebra"
    Animal <|-- Duck
    note for Duck "can fly\ncan swim\ncan dive\ncan help in debugging"
    Animal <|-- Fish
    Animal <|-- Zebra
    Animal : +int age
    Animal : +String gender
    Animal: +isMammal()
    Animal: +mate()
    class Duck{
        +String beakColor
        +swim()
        +quack()
    }
    class Fish{
        -int sizeInFeet
        -canEat()
    }
    class Zebra{
        +bool is_wild
        +run()
    }
```
