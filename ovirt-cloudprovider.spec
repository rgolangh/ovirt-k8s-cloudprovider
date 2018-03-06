Name:       ovirt-cloudprovider
Version:    %{?_version}
Release:    %{?_release}%{?dist}
Summary:    oVirt external cloud-provider for kubernetes/OpenShift

License:    ASL 2.0
URL:        http://www.ovirt.org
Source0:    %{name}-%{version}%{?_release:-%_release}.tar.gz

%description
oVirt external cloud-provider for kubernetes/OpenShift

%prep
%setup -c


%build
echo $(pwd)
go env
make deps
make build

%install
install -p -m 755 %{name} %{buildroot}/bin/

%define debug_package %{nil}

%files
%{buildroot}/bin/%{name}

%changelog