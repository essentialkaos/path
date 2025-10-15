################################################################################

%global crc_check pushd ../SOURCES ; sha512sum -c %{SOURCE100} ; popd

################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Dead simple tool for working with paths
Name:           path
Version:        1.2.0
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/path

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

Source100:      checksum.sha512

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.24

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Dead simple tool for working with paths.

################################################################################

%prep
%{crc_check}

%setup -q
if [[ ! -d "%{name}/vendor" ]] ; then
  echo -e "----\nThis package requires vendored dependencies\n----"
  exit 1
elif [[ -f "%{name}/%{name}" ]] ; then
  echo -e "----\nSources must not contain precompiled binaries\n----"
  exit 1
fi

%build
pushd %{name}
  %{__make} %{?_smp_mflags} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_mandir}/man1

install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

./%{name}/%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

%clean
rm -rf %{buildroot}

%post
if [[ -d %{_sysconfdir}/bash_completion.d ]] ; then
  %{name} --completion=bash 1> %{_sysconfdir}/bash_completion.d/%{name} 2>/dev/null
fi

if [[ -d %{_datarootdir}/fish/vendor_completions.d ]] ; then
  %{name} --completion=fish 1> %{_datarootdir}/fish/vendor_completions.d/%{name}.fish 2>/dev/null
fi

if [[ -d %{_datadir}/zsh/site-functions ]] ; then
  %{name} --completion=zsh 1> %{_datadir}/zsh/site-functions/_%{name} 2>/dev/null
fi

%postun
if [[ $1 == 0 ]] ; then
  if [[ -f %{_sysconfdir}/bash_completion.d/%{name} ]] ; then
    rm -f %{_sysconfdir}/bash_completion.d/%{name} &>/dev/null || :
  fi

  if [[ -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish ]] ; then
    rm -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish &>/dev/null || :
  fi

  if [[ -f %{_datadir}/zsh/site-functions/_%{name} ]] ; then
    rm -f %{_datadir}/zsh/site-functions/_%{name} &>/dev/null || :
  fi
fi

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_mandir}/man1/%{name}.1.*
%{_bindir}/%{name}

################################################################################

%changelog
* Wed Oct 15 2025 Anton Novojilov <andy@essentialkaos.com> - 1.2.0-0
- Dependencies update

* Thu May 08 2025 Anton Novojilov <andy@essentialkaos.com> - 1.1.1-0
- Added info about supported environment variables to usage info
- Minor UI improvements
- Code refactoring
- Dependencies update

* Fri Nov 08 2024 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Added 'strip-ext' command
- Code refactoring
- Dependencies update

* Tue Sep 24 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.3-0
- Dependencies update

* Sun Jun 23 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.2-0
- Code refactoring
- Dependencies update

* Thu Mar 28 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Improved support information gathering
- Code refactoring
- Dependencies update

* Sun Feb 18 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Added dirn command
- Code refactoring
- Dependencies update

* Sun Dec 17 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.6-0
- Improved verbose version info output
- Code refactoring
- Dependencies update

* Tue May 23 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.5-0
- Add support of enabling quiet mode using environment variable (PATH_QUIET)

* Wed May 17 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.4-0
- Improve input parsing

* Tue May 16 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.3-0
- Added add-prefix command
- Added remove-prefix command
- Added add-suffix command
- Added remove-suffix command
- Added exclude command

* Mon May 15 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.2-0
- Added join command
- Custom version info formats support
- Added bibop tests
- Code refactoring

* Thu May 04 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- First public release
