# <div id="top"></div>

<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![Wiki][wiki-shield]][wiki-url]

</div>

  <h1 align="center">ERP Login Module</h1>

  <p align="center">
  <!-- UPDATE -->
    <a href="https://github.com/metakgp/iitkgp-erp-login-go/issues">Report Bug</a>
    ·
    <a href="https://github.com/metakgp/iitkgp-erp-login-go/issues">Request Feature</a>
  </p>
</div>


<!-- TABLE OF CONTENTS -->
<details>
<summary>Table of Contents</summary>

- [About The Project](#about-the-project)
  - [Supports](#supports)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Projects using ERP login pacakge Go](#projects-using-erp-login-pacakge-go)
- [Maintainer(s)](#maintainers)
- [Contact](#contact)
- [Additional documentation](#additional-documentation)

</details>


## About The Project

This package automates the login workflow for IIT Kharagpur ERP.

<p align="right">(<a href="#top">back to top</a>)</p>

<div id="supports"></div>

### Supports:
1. Shells
    * `bash`
    * `zsh`
2. OS(s)
    * any `*nix`[`GNU+Linux` and `Unix`]

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

To set up a local instance of the application, follow the steps below.

### Prerequisites
The following dependencies are required to be installed for the project to function properly:
<!-- UPDATE -->
- [Go](https://go.dev/)
- IIT-KGP student ERP account

<p align="right">(<a href="#top">back to top</a>)</p>

### Installation

_Now that the environment has been set up and configured to properly compile and run the project, the next step is to install and configure the project locally on your system._
<!-- UPDATE -->
1. Clone the repository
   ```sh
   git clone https://github.com/metakgp/iitkgp-erp-login-go
   ```
2. Install dependencies
   ```sh
   cd ./iitkgp-erp-login-go
   go mod download
   ```
3. Create a `erpcreds.json` with following contents for auto-login (optional)
    ```json
      {
        "roll_number": "Enter roll number",
        "password": "Enter ERP password",
        "answers": {
          "Security Question 1": "Answer to security 1",
          "Security Question 2": "Answer to security 2",
          "Security Question 3": "Answer to security 3"
        }
      }
    ```


<p align="right">(<a href="#top">back to top</a>)</p>

## Usage
The following code logs in to ERP and opens ERP homepage in the browser.

```go
package main

import (
	erp "github.com/metakgp/iitkgp-erp-login-go"

	"github.com/pkg/browser"
)

func main() {
	_, ssoToken := erp.ERPSession()
	
	browser.OpenURL(erp.HOMEPAGE_URL + "?" + ssoToken)
}
```

### Projects using ERP Login Package Go
- [Chillzone](https://github.com/metakgp/chillzone)

## Maintainer(s)

- [Shikhar Soni](https://github.com/shikharish)

<p align="right">(<a href="#top">back to top</a>)</p>

## Contact

<p>
📫 MetaKGP -
<a href="https://slack.metakgp.org">
  <img align="center" alt="Metakgp's slack invite" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/slack.svg" />
</a>
<a href="mailto:metakgp@gmail.com">
  <img align="center" alt="Metakgp's email " width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/gmail.svg" />
</a>
<a href="https://www.facebook.com/metakgp">
  <img align="center" alt="metakgp's Facebook" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/facebook.svg" />
</a>
<a href="https://www.linkedin.com/company/metakgp-org/">
  <img align="center" alt="metakgp's LinkedIn" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/linkedin.svg" />
</a>
<a href="https://twitter.com/metakgp">
  <img align="center" alt="metakgp's Twitter " width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/twitter.svg" />
</a>
<a href="https://www.instagram.com/metakgp_/">
  <img align="center" alt="metakgp's Instagram" width="22px" src="https://raw.githubusercontent.com/edent/SuperTinyIcons/master/images/svg/instagram.svg" />
</a>
</p>

<p align="right">(<a href="#top">back to top</a>)</p>

## Additional documentation

  - [License](/LICENSE)
  - [Code of Conduct](/.github/CODE_OF_CONDUCT.md)
  - [Security Policy](/.github/SECURITY.md)
  - [Contribution Guidelines](/.github/CONTRIBUTING.md)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->

[contributors-shield]: https://img.shields.io/github/contributors/metakgp/iitkgp-erp-login-go.svg?style=for-the-badge
[contributors-url]: https://github.com/metakgp/iitkgp-erp-login-go/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/metakgp/iitkgp-erp-login-go.svg?style=for-the-badge
[forks-url]: https://github.com/metakgp/iitkgp-erp-login-go/network/members
[stars-shield]: https://img.shields.io/github/stars/metakgp/iitkgp-erp-login-go.svg?style=for-the-badge
[stars-url]: https://github.com/metakgp/iitkgp-erp-login-go/stargazers
[issues-shield]: https://img.shields.io/github/issues/metakgp/iitkgp-erp-login-go.svg?style=for-the-badge
[issues-url]: https://github.com/metakgp/iitkgp-erp-login-go/issues
[license-shield]: https://img.shields.io/github/license/metakgp/iitkgp-erp-login-go.svg?style=for-the-badge
[license-url]: https://github.com/metakgp/iitkgp-erp-login-go/blob/master/LICENSE
[wiki-shield]: https://custom-icon-badges.demolab.com/badge/metakgp_wiki-grey?logo=metakgp_logo&style=for-the-badge
[wiki-url]: https://wiki.metakgp.org