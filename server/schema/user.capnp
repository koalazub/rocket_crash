using Go = import "/go.capnp";

@0xe26cbc744ed2a7c9;

$Go.package("schema");
$Go.import("schema/user.capnp");

struct User {
  user @0 :Text;
  points @1 :Int64;
}
