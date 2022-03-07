### Development
We have built an AMI using packer and you can find the code for in the packer directory.
This code uses packer version `1.8.0`
```
packer init .
packer validate .
packer build .
```

Using that AMI we have started an instance in us-west-2 in the default VPC's public subnet where you can run the build.
The code for that is contained in the terraform directory
```
terraform init
terraform plan
terraform apply
```

The private key for connecting to the instance is in an AWS SSM parameter named `sm`. You can grab the RSA private key from there to ssh into the machine.
Username for ssh is `ubuntu`

### Build
How to use the ec2 instance to build a package
```
sudo su builder
cd ~
cd build

git clone https://github.com/brave/brave-browser.git
cd brave-browser
npm install
npm run init
./src/build/install-build-deps.sh
npm run build
```
