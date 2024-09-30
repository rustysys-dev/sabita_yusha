Summary: keyboard macro configuration and execution daemon
Name: sabita_yusha
Version: 1.0.0
Release: 1
License: GPLv3+
Packager: rustysys-dev

%description
A configurable macro keyboard key-press executor for programmable-button 
1 - 30 configured using VIA/QMK/REMAP.  

%prep
mkdir -p %{buildroot}/src
cp -r * %{buildroot}/src/

%build
cd %{buildroot}/src
go build main.go -o sabita_yusha

%install
mkdir -p %{buildroot}%{_bindir}
install -m 0755 sabita_yusha %{buildroot}%{_bindir}/sabita_yusha
mkdir -p %{buildroot}%{_userunitdir}
install -m 0644 scripts/systemd/sabita_yusha.service %{buildroot}%{_userunitdir}/sabita_yusha.service

%post
systemctl --user enable sabita_yusha.service
systemctl --user start sabita_yusha.service

%files
# _bindir resolves to /usr/bin 
# https://docs.fedoraproject.org/en-US/packaging-guidelines/RPMMacros/#macros_installation
%{_bindir}/sabita_yusha
# _userunitdir resolves to /usr/lib/systemd/user 
# https://docs.fedoraproject.org/en-US/packaging-guidelines/Systemd/#packaging
%{_userunitdir}/sabita_yusha.service

%changelog
* Mon Sep 30 2024 rustysys-dev <scott.mattan@rustysys.dev>
- initial build
