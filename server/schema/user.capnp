using Go = import "/go.capnp";

@0x9c76bc62a71f9389;
$Go.package("user");
$Go.import("schema/user.capnp");

struct User {
  user @0 :Text;
  points @1 :Int64;
}
