using Go = import "/go.capnp"; # I think this is the import statement for the go version of capnp
@0x9c76bc62a71f9389;
$Go.package("schema");
$Go.import("schema/rocket.capnp");

struct Rocket {
  name @0 :Text;
  user @1 :Text;  
  deathCoordY @2 :Float64;
  deathCoordX @3 :Float64;
  rocketType @4 :Text;
  crashed @5 :Bool;
}  
