# Download Manager Golang

Welcome to the Download Manager Golang repository. This project is a download manager implemented in Go (Golang) designed to efficiently handle and manage multiple file downloads.

## Repository
https://github.com/SUT-technology/download-manager-golang

## Features

- Concurrent Downloads: Supports downloading multiple files simultaneously.
- Pauseable  &Resumable Downloads: Allows pausing and resuming downloads.
- Organising Downloads: You can download in multiple queues and schedule your downloads in specific times of the day
- Progress Tracking: Displays download progress for each file.

## Getting Started

### Prerequisites

- Go 1.16+ installed on your machine. You can download and install it from [here](https://golang.org/dl/).

### Installation

1. Clone the repository:

   
   git clone https://github.com/SUT-technology/download-manager-golang.git
   cd download-manager-golang
   

2. Build the project:

   
   go build
   

### Usage

To start the download manager, run the following command:


./download-manager-golang


## Project Description

The clean design architecture is used for the file structure of this project that contains four folders. First we have the assets folder, which contains the database (two json files for saving download and queues) and the project config. Then there is the cmd package which contains the runnable files (main.go and run.go). The pkg folder contains some tools used in the project. Last but not least, we have the internal folder that is the main part of the project. It contains the user interface of the program which is written by the BubbleTea library, the download and queue entities and dtos in the domain package, some interfaces for download and queue logic, the pool.go file in the infrastructure folder that connects the program with the json database, the services that contain the main logic (e.g. downloading files), the handlers that connect the ui to the service layer, and etc. 



## Contributors 

Mahdi Rasoulzadeh (402170121, https://github.com/its-mahdi-dev) Created the initial structure of the project and wrote the concurrency and pause & resume logic.
Souroosh Najafi (402100559, https://github.com/CDsonji) Wrote the structure of the ui and queue editing and deleting logic and ui.
Erfan Ghorbani (402170216, https://github.com/erfan23g) Wrote the logic of downloading and creating queues and the ui of adding downloads.