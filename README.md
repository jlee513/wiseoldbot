<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Issues][issues-shield]][issues-url]
[![Apache License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="images/wiseoldman.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">Wise Old Bot</h3>

  <p align="center">
    A fully fledged Old School Runescape Discord bot that allows for the entire clan experience!
    <br />
    <a href="https://github.com/jlee513/wiseoldbot"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/jlee513/wiseoldbot">View Demo</a>
    ·
    <a href="https://github.com/jlee513/wiseoldbot/issues">Report Bug</a>
    ·
    <a href="https://github.com/jlee513/wiseoldbot/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#setup-information">Setup Information</a>
      <ul>
        <li><a href="#admin-discord-channel">Admin Discord Channel</a></li>
        <li><a href="#discord-channels-for-the-bot">Discord Channels for the Bot</a></li>
        <li><a href="#how-to-setup-pastebin-for-guides">How to setup Pastebin for Guides</a></li>
      </ul>
    </li>
    <li><a href="#general-information">General Information</a>
      <ul>
        <li><a href="#main-player-paradigm">Main Player Paradigm</a></li>
      </ul>
    </li>
    <li><a href="#clan-points-instructions">Clan Points Instructions</a>
      <ul>
        <li><a href="#player-submission">Player Submission</a></li>
        <li><a href="#admin-approvals">Admin Approvals</a></li>
        <li><a href="#feedback-channel">Feedback Channel</a></li>
      </ul>
    </li>
    <li><a href="#speed-instructions">Speed Instructions</a>
      <ul>
        <li><a href="#player-submission">Player Submission</a></li>
        <li><a href="#admin-approvals">Admin Approvals</a></li>
        <li><a href="#feedback-channel">Feedback Channel</a></li>
      </ul>
    </li>
    <li><a href="#guide-slash-command-instructions">Guide Slash Command Instructions</a>
      <ul>
        <li><a href="#guide">/guide</a></li>
      </ul>
    </li>
   <li><a href="#admin-slash-command-instructions">Admin Slash Command Instructions</a>
      <ul>
        <li><a href="#admin-player">/admin player</a></li>
        <li><a href="#admin-speed">/admin speed</a></li>
        <li><a href="#admin-leaderboard">/admin leaderboard</a></li>
        <li><a href="#admin-sheets">/admin sheets</a></li>
        <li><a href="#admin-instructions">/admin instructions</a></li>
        <li><a href="#admin-guides-map">/admin guides-map</a></li>
        <li><a href="#admin-points">/admin points</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

![Product Name Screen Shot][product-screenshot]

Every Old School Runescape clan needs a way to keep track of members' progress that feel meaningful. It always is a pain to go to another website just to check who get new achievements or to manually keep track of clan points via a ticketing system that feels archaic. That's where Wise Old Bot comes in.

Here's an overview of features available:
* Built in clan point management and speed tracking
* Discord slash commands for clan members that allow for clan point & speed submissions
* Admin slash commands for clan administration
* Guide slash commands for guide addition/modification
* Integration with third-party Runescape metrics websites such as TempleOSRS and Collectionlog.net
* Leaderboards/Hall Of Fame based off of these metrics
* Automatic Clan Point addition based off the Better Discord Loot Logger plugin on Runelite
* Feedback Discord Channel that populates when a decision has been made about a submission

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

The programming language used to build Wise Old Bot is Golang. It leverages a well documented and robust library from bwmarrin called discordgo.

* [![Go][Golang]][Golang-url]
* [![discordgo](https://pkg.go.dev/badge/github.com/bwmarrin/discordgo.svg)](https://pkg.go.dev/github.com/bwmarrin/discordgo)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

You need to have Golang installed on your computer. For instructions on how to download, click the following link to go to the go download and installation instructions:

https://go.dev/doc/install

You will also need the accounts on the following third party applications:

* Discord Bot created for this project
* Google account that all the admins will have access to
* Pastebin account that all the guide admins will have access to
* TempleOSRS group that all the members will be added to
* Imgur account that all non-imgur user submissions will be uploaded to
* Postman for Imgur Setup and ease of API testing

### Installation

_Below is an example of how you can instruct your audience on installing and setting up your app. This template doesn't rely on any external dependencies or services._

1. Clone the repo
   ```sh
   git clone https://github.com/jlee513/wiseoldbot.git
   ```
2. Create a Discord Bot Account that will be used to handle all of your requests. The following URL describes how to do that step by step.
   ```txt
   https://discordpy.readthedocs.io/en/stable/discord.html
   ```
3. You want to create an individual Google Sheets for each of the following purposes. Each of your sheets must have their share settings set to Anyone With The Link as an editor. Each of the sheets need to follow a particular structure and header as defined under each of the sheets.
   ```txt
   Members:
   Player | Discord ID | Discord Name | Feedback | Main
   
   Clan Points:
   Player | Clan Points
   
   Clan Points Submission Screenshots:
   Time Submitted | URL | Player
   
   Speed Times:
   Category | Boss Name | Time | PLayers | URL
   
   Speed Times Submission Screenshots:
   Time Submitted | URL | Boss Name | Time | PLayer
   
   Transaction ID:
   Nothing is needed for the transaction ID sheets page.
   It will populate cell A1.
   
   ```
4. Register an Imgur Application following the steps outlined in the Imgur API docs
   ```txt
   https://apidocs.imgur.com/
   ```
5. Create client_secret.json file under the config folder of this project. Follow the instructions in the following URL:
   ```txt
   https://help.talend.com/r/en-US/8.0/google-drive/oauth-methods-for-accessing-google-drive
   ```
6. Create a .env file under the config folder of this project. Look under example/.env_test for an outline of all the current environment variables that are required to run the bot.
   ```txt
   Temple Information:
   
   TEMPLE_GROUP_ID=The Unique identifier for your clan
   TEMPLE_GROUP_KEY=The main key used to update clan members
   ```
   ```txt
   Google Sheets:
   
   SHEETS_CP=Sheets ID for clan points
   SHEETS_CP_SC=Sheets ID for clan points screenshots
   SHEETS_SPEED=Sheets ID for speed times
   SHEETS_SPEED_SC=Sheets ID for speed times screenshots
   SHEETS_TID=Sheets ID for transaction id
   SHEETS_MEMBERS=Sheets ID for members
   ```
   ```txt
   Imgur Information:
   
   IMGUR_CLIENT_ID=ID of your imgur client
   IMGUR_CLIENT_SECRET=Secret of your imgur client
   IMGUR_REFRESH_TOKEN=Refresh Token used to get access tokens for uploaded imgurs
   ```
   ```txt
   Pastebin:
   
   PASTEBIN_USERNAME=Username of your Pastebin account
   PASTEBIN_PASSWORD=Password of your Pastebin account
   PASTEBIN_DEV_API_KEY=To get your dev api key, login to Pastebin and navigate to the API tab. Here, you can find a section called "Your Unique Developer API Key"
   PASTEBIN_MAIN_PASTE_KEY=Pastebin key id that will hold a list of your guides
   ```
   ```txt
   Discord Channels:
   
   DISCORD_*_CHANNEL=Discord Channel ID
   ```
5. Compile and run the bot
   ```sh
   go run .
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Setup Information

### Admin Discord Channel

When you add the bot, ensure that the admin and guide channels are locked behind roles and locked behind specific discord channels so that not everyone is able to use them. 

![Bot Integration][bot-integration-screenshot]

### Discord Channels for the Bot

When you have a community discord, you have the ability to create Forums which is a great way to condense sections for leaderboards, hall of fame for kcs and speed, as well as guides. 

Whenever the bot posts to a discord channel, it removes all messages from the channel before posting. To ensure that the picture in Forum posts do not get deleted, make sure to use the picture upon creation of the post. This allows the bot to know that it should delete all messages after the picture.

![Guides Forum][guides-screenshot]

### How to setup Pastebin for Guides

Create one pastebin that will serve as the holder of the keys as well as the discord channel it will populate. The following will illustrate what this file should look like.
   ```txt
   # Tob
   Tob Fives - xxxxxxxx - 8648323223567591832
   Tob Fives - xxxxxxxx - 8648323223567591832
   Tob Fives - xxxxxxxx - 8648323223567591832
   Tob Fives - xxxxxxxx - 8648323223567591832
   # CM Trio
   CM Trio Prep - xxxxxxxx - 8648323223567591832
   CM Trio Chin - xxxxxxxx - 8648323223567591832
   CM Trio Surge - xxxxxxxx - 8648323223567591832 
   CM Trio Useful Info - xxxxxxxx - 8648323223567591832
   ```
Help: The # denotes the overall guide and each line under it will be it's own unique discord channel. Within these lines, the first will be a human readable name that isn't used in the creation of the guide - it's more for you to know what this page is for. The x's is the pastebin key which you can find in the URL when you go to a pastebin. The last number is the channel ID that you want the guide to populate in.

_For information about guides, please refer to the [Wise Old Bot Guide Creation Instruction](https://pastebin.com/aas7Lets)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GENERAL INFORMATION -->
## General Information

### Main Player Paradigm
When it comes to players, we know that many people can have an account that they mainly play but also have alts. There is a feature that allows you to combine the KCs & Clan Points under one main username so that it tracks you as a player instead of individual accounts separately.

![Main][main-screenshot]

In this example above, the clan points and speed times will be tracked under the player "Mager". However, since "sry lil bro" is in the clan but is not the main, all speed times/clan points/kcs will be added and displayed under "Mager".

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CP AND SPEED SLASH COMMAND INSTRUCTIONS -->
## Clan Points Instructions
### Player Submission
![cpsubmit][cpsubmit-gif]

### Admin Approvals
Once a successful submission is submitted as shown above, a message will pop up in a designated hidden channel meant for only moderators. You can approved/denied using the ✅ or ❌ buttons at the bottom. Once selected, the message will be removed from the channel and will send feedback to the appropriate feedback channel of the submitter.

![cpadminapproval][cpadminapproval-screenshot]

### Feedback Channel
A feedback channel will be created if it's the first time someone has submitted. It will create a channel with only that person as well as the moderators. This is the place that is used to discuss why a submission was accepted/rejected as well as serves a history of all submissions for the user.

![feedbackcp][feedbackcp-screenshot]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Speed Instructions
### Player Submission
![speedsubmit][speedsubmit-gif]

### Admin Approvals
Once a successful submission is submitted as shown above, a message will pop up in a designated hidden channel meant for only moderators. You can approved/denied using the ✅ or ❌ buttons at the bottom. Once selected, the message will be removed from the channel and will send feedback to the appropriate feedback channel of the submitter.

![speedapproval][speedapproval-screenshot]

### Feedback Channel
A feedback channel will be created if it's the first time someone has submitted. It will create a channel with only that person as well as the moderators. This is the place that is used to discuss why a submission was accepted/rejected as well as serves a history of all submissions for the user.

![feedbackspeed][feedbackspeed-screenshot]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GUIDE SLASH COMMAND INSTRUCTIONS -->
## Guide Slash Command Instructions
![guidecommands][guidecommands-screenshot]

### /guide
Guide Updates Command
```txt
- Update: Updates the guide selected
  Required Options: option, guide
```
Note: This autopopulated guide list is updated by the /admin guides-map command based off of the main pastebin.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ADMIN SLASH COMMAND INSTRUCTIONS -->
## Admin Slash Command Instructions
![admincommands][admincommands-screenshot]

### /admin player
Player administration commands
```txt
- Add: Add player to the clan (automatically adds to temple)
  Required Options: option, name, discord-id, discord-name, main
  
- Remove: Remove player from clan (automatically removes from temple)
  Required Options: option, name

- Name Change: Moves player information from one name to another
  Required Options: option, name, new-name
  
- Update Main: Updates the main username tracking for a player
  Required Options: option, name, new-name
```

### /admin speed
Speed administration commands
```txt
- Add: Adds a new speed time under an existing category
  Required Options: option, category, new-boss
  
- Remove: Removed a speed time under an existing category
  Required Options: option, category, existing-boss

- Update: Updates the speed time, player names, or the image for an existing boss
  Required Options: option, category, existing-boss
  Optional Options: update-speed-time, update-player-names, imgur-url
  
- Reset: Resets the speed time for an existing boss
  Required Options: option, category, existing-boss
```
Note: category and existing-boss autopopulates based on the speed sheets list.

### /admin leaderboard
Leaderboard administration commands
```txt
- Update: Updates the leaderboard selected
  Required Options: option, leaderboard, thread
```

Note: leaderboard and thread autopopulates which is hardcoded.

### /admin sheets
Google Sheets administration commands
```txt
- Update: Updates the google sheet based on the maps/arrays stored in memory
  Required Options: option
```

### /admin instructions
Update clan points and speed submission instructions
```txt
- Update: Update the clan points and speed submission instructions from code
  Required Options: option
```

### /admin guides-map
Update map of guides from the main pastebin
```txt
- Update: Update map of guides from the main pastebin
  Required Options: option
```

### /admin points
Update points of the comma separated list of players
```txt
- Add: Adds the specified number of clan points
  Required Options: option, player, amount-of-pp
  
- Remove: Removes the specified number of clan points
  Required Options: option, player, amount-of-pp
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

See the [trello board](https://trello.com/b/84ZHKrbf/wise-old-bot) for a roadmap on what's to do and what's in progress.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Joshua Lee - [@MagerOsrs](https://twitter.com/MagerOsrs)

Project Link: [https://github.com/jlee513/wiseoldbot](https://github.com/jlee513/wiseoldbot)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/jlee513/wiseoldbot.svg?style=for-the-badge
[contributors-url]: https://github.com/jlee513/osrs-disc-bot/graphs/contributors
[issues-shield]: https://img.shields.io/github/issues/jlee513/wiseoldbot.svg?style=for-the-badge
[issues-url]: https://github.com/jlee513/wiseoldbot/issues
[license-shield]: https://img.shields.io/github/license/jlee513/wiseoldbot.svg?style=for-the-badge
[license-url]: https://github.com/jlee513/osrs-disc-bot/blob/main/LICENSE
[product-screenshot]: images/banner.png
[guides-screenshot]: images/guidechannels.png
[bot-integration-screenshot]: images/botintegrations.png
[main-screenshot]: images/main.png
[admincommands-screenshot]: images/admincommands.png
[guidecommands-screenshot]: images/guidecommands.png
[cpsubmit-gif]:images/cpsubmit.gif
[cpadminapproval-screenshot]: images/cpadminapproval.jpg
[feedbackcp-screenshot]: images/feedbackcp.png
[speedsubmit-gif]: images/speedsubmit.gif
[speedapproval-screenshot]: images/speedapproval.jpg
[feedbackspeed-screenshot]: images/feedbackspeed.png
[Golang]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[Golang-url]: https://go.dev/

