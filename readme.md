# Tanks-On-Sale-Notifier
A crawler to fetch latest vehciles on sale from [360 WOT Shop](https://shop.wot.360.cn) and post it on twitter.


# Why brother?
Wot shop sells golden vehciles only for a brief limited time. If you only want a specific vehcile and do not have time logged in the game everyday, you'll surely miss the sale and have to wait for another few months.


# How it works:
The crawler is written in go, and invocated by github actions every hour. It will fetch the shop, distill the info to concise form, and broadcast it on twitter/telegram. 


# Example output:
```
TS-5          ￥227.50  至11-01 10:00
703工程 II型  ￥265.00  至10-29 10:00
```


# To-Do:
- [ ] Apply for Twitter API / Telegram API 
- [ ] Logic for posting
- [ ] Unit Test for functions