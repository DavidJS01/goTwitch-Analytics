# Go-Twitch Analytics
<p float="left"> <b>Golang meets Twitch chat!</b> </p>

<p float="left" align="">
  <img align = center src="assets/readme/gopher-dance-long-3x.gif"  />
  <img align=center src="assets/readme/transparent_twitch.png" width=200px /> 
</p>


## What is Go-Twitch Analytics?
Go-Twitch is a self-service, containerized solution to aggregating and visualizing Twitch chat data in near-real time.

Go-Twitch is a rewrite of my first data engineering Python project, but written in Go (with more elaborate technology).

Users can go into the web-app (WIP) and select their own list of Twitch chat streams to aggregate those messages before having those messages analyzed, visualizing the data in near real-time.

This is in effort to continously learn and exit my comfort zones to become a great data engineer.
## Example Dashboard
![](assets/readme/example_dashboard.png)

## Demo
[![Go-Twitch Analytics](assets/readme/Go-Twitch%20Thumbnail.png)](https://www.youtube.com/watch?v=wSOO38p4rNw "Go-Twitch Analytics Demo")
<p align="center">(redirects to youtube)</p>

## System Architecture Diagram
![](assets/readme/GoTwitchV2.png)

## Entity Relationship Diagram
![](assets/readme/gotwitch_erd.drawio.png)   


## Disclaimers and Notes
Development on GoTwitch is currently on hold.

This project is continously being updated and developed, with architecture and system design improvements being prioritized. *It is not yet fit for consumer use yet*.

Due to budget constraints, this project is locally containerized. In the ideal world, it would be hosted on AWS and be a public facing website. The data architecture would probably just be an RDS instance. 

Top priority development goals:
1. Redesigned, normalized data models
2. OAuth handled programatically
3. MVP front-end UI for a user-friendly, accessible experience

![](assets/readme/goTwitchGanttChart.png)

