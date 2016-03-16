// START OMIT
#Defining the provider
provider "aws" {
  region = "${var.aws_region}"
}

resource "aws_eip" "nat" {
  vpc = true
}

resource "aws_nat_gateway" "gw" {
    allocation_id = "${aws_eip.nat.id}"
    subnet_id = "${lookup(var.aws_public_subnets, var.aws_region)}"
}

#Create a private subnet for the instances
resource "aws_subnet" "private" {
  vpc_id      = "${lookup(var.aws_vpc, var.aws_region)}"
  cidr_block = "172.31.32.0/20"
  tags {
      Name = "Private Subnet"
  }
}
// END OMIT
