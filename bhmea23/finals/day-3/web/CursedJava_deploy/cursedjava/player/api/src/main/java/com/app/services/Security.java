 package com.app.services;
 import java.io.*;


 import java.util.regex.Pattern;
 
 public class Security extends ObjectInputStream {
   public Security(InputStream inputStream) throws IOException {
     super(inputStream);
   }
 
   
   protected Class<?> resolveClass(ObjectStreamClass desc) throws IOException, ClassNotFoundException {
     if (!Pattern.matches("(com\\.app\\.(.*))|(java\\.time\\.(.*))", desc.getName())) {
       throw new InvalidClassException("Unauthorized deserialization attempt", desc.getName());
     }
     return super.resolveClass(desc);
   }
 }
