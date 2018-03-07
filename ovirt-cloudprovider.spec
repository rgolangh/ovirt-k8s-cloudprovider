Name:       ovirt-cloudprovider
Version:    %{?_version}
Release:    %{?_release}%{?dist}
Summary:    oVirt external cloud-provider for kubernetes/OpenShift

License:    ASL 2.0
URL:        http://www.ovirt.org
Source0:    %{name}-%{version}%{?_release:_%_release}.tar.gz

%global repo github.com/rgolangh
%global golang_version 1.9.1
%global debug_package %{nil}

%description
oVirt external cloud-provider for kubernetes/OpenShift

%prep
%setup -c -q

%build
# set up temporary build gopath for the rpmbuild
mkdir -p ./_build/src/%{repo}
ln -s $(pwd) ./_build/src/%{repo}/%{name}

export GOPATH=$(pwd)/_build
cd _build/src/%{repo}/%{name}
go env
make build

%install
mkdir -p %{buildroot}/usr/bin
install -p -m 755 %{name} %{buildroot}/usr/bin/

%files
/usr/bin/%{name}

%changelog