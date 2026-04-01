resource "aws_security_group" "golang_sg" {
  name        = "${var.name_project}-sg"
  description = "Libera o acesso do EKS"
  vpc_id      = aws_vpc.golang_vpc.id

  ingress {
    description = "Acesso do RebbitMQ"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.aws_region]

  }



  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [0.0.0.0/0]
  }
}