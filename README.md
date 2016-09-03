![image](./resources/frankinstore-logo.jpg)

##about

Frankinstore is a basic persistent, content addressable, blob store, backed by bolted DB. Binary blobs can be of arbitrary non-zero size. Keys are the 20 byte computed SHA-1 signature of the blob. 

Frankinstore exposes 2 web endpoints for the supported `Get` and `Put`

## instalation

Clone the project and run the top-level server.go. 

    git clone https://github.com/alphazero/frankinstore.git
    cd frankinstore
    
    # run using -h to get options details
    go run server.go [options]
    

## API


###Put

Put is a POST method call to the service. If successful (http-stat 200), the response body is the associated key of the blob. Note that the key is returned as binary (e.g. 20 bytes) and is not encoded.
 
     method:    POST
     uri:       /put
     body:      <binary blob>
  
example (assuming localhost:5722):

     http://localhost:5722/put

###Get

Get is a simple GET method call to the service. If successful (http-stat 200), the response body is the value binary blob.
 
     method:    GET
     uri:       /get/<hex-encoded-key>
  
example (assuming localhost:5722):

     http://localhost:5722/get/316eb0ec4c0f75f4cbb19b6b5e59142e0fb01214
     
## server options

You can specifiy the location/name of the FS store, and of course the port for the service apis.


## NOTICE 3rd Party Software

Frankinstore uses `bolt/boltdb` and a specific package from `google/groupcache`. See the included `LICENSE-3rdParty` file for details.

## License 

    Copyright © 2016 Joubin Houshyar. All rights reserved.

    This file is part of Frankinstore.

    Frankinstore is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of
    the License, or (at your option) any later version.

    Frankinstore is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public
    License along with Frankinstore.  If not, see <http://www.gnu.org/licenses/>.
