# pong
Make the old-time classic game 'pong' using GameSalad


## Resolution
480x320

Fixing paddle bounds (measure from the middle of the paddle)
The paddle has: 16x64, so Y/2 = 32

* Minimum Y for the paddle should be: 0 + 32 = 32  ==>  Y  >= 32
* Maximym Y for the paddle should be: 320 - 32 = 288  ==> Y <= 288 


