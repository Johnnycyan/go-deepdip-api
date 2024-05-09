<p align="center">
  <img src="https://raw.githubusercontent.com/Johnnycyan/Twitch-APIs/main/OneMoreDayIcon.svg" width="100" alt="project-logo">
</p>
<p align="center">
    <h1 align="center">GO-DEEPDIP-API</h1>
</p>
<p align="center">
    <em>Dive Deep, Rise High with Real-Time Insights</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/last-commit/Johnnycyan/go-deepdip-api?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/Johnnycyan/go-deepdip-api?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/Johnnycyan/go-deepdip-api?style=default&color=0080ff" alt="repo-language-count">
<p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>

<hr>

##  Overview

The go-deepdip-api is a software project designed to enhance user experience in Trackmania's Deep Dip event by providing real-time updates. It establishes a server that fetches player details and leaderboard rankings, significantly improving dynamic interactions through its integration with external APIs. This is further supported by its well-managed dependencies and automated workflows, which guarantee the consistency and reliability of the deployment. The use of environmental configurations and external libraries allows the go-deepdip-api to remain flexible and efficient, making it a valuable tool for any Trackmania degen.

---

##  Getting Started

**System Requirements:**

* **Internet**

###  Installation

<h4>From <code>releases</code></h4>

> 1. Download the latest release:
>     1. [Latest Releases](https://github.com/Johnnycyan/go-deepdip-api/releases) 
>
> 2. Create .env file in the same directory and add `NAME=<your-username>`

###  Usage

<h4>From <code>releases</code></h4>

> Run go-deepdip-api using the command below:
> ```console
> $ ./go-deepdip-api <port>
> ```
> or
> ```console
> $ go-deepdip-api.exe <port>
> ```
Endpoint      |     URL
------------- | -------------
PB  | `localhost:<port>/pb?username=<player-name>`
Leaderboards  |  `localhost:<port>/leaderboards?username=<player-name>` (username optional)

---
