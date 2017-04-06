## yaml2p 

Simplistic converter between yaml and java properties files.

### Not supported

* multiline string values
* yaml collection of the type `key: [val1, val2, val3]`
* scalars at a node with children, e.g. the following java property file won't be mapped to yaml: 
  
````
level1.level2=blah
level1.level2.level3=fooe
````  

### Quirks

Original key order isn't preserved, alphabetical sorting is used for the consistency sake