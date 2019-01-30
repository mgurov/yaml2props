## yaml2p 

Simplistic converter between yaml and java properties files.

## Installation 

Download the latest binary from the [github releases archive](https://github.com/mgurov/yaml2props/releases)

## from sources (GO)

````bash
go get github.com/mgurov/yaml2props/cmd/yaml2p
go install github.com/mgurov/yaml2props/cmd/yaml2p
````

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

## Related art 

Java code https://stackoverflow.com/questions/49207935/how-to-convert-yml-to-properties-with-a-gradle-task 

Java/maven style https://github.com/redlab/yaml-props
