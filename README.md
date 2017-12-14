```
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
```


## Install

`go get -u github.com/bigkucha/go-socket`

## Demo

Run there terminal windows, One for server, tow for clients.

the first window for client ,userID: 1 , userIDWhoYouWantToTalk: 2  
the second window for client , userID: 2, userIDWhoYouWantToTalk: 3

type message in the first client window and click enter button, and the second window will received!

```
$ cd $GOPATH/src/github/bigkucha/go-socket/socket-demo
$ go build
$ ./socket-demo server
$ ./socket-demo client {userID} {userIDWhoYouWantToTalk}
$ ./socket-demo client {UserID} {userIDWhoYouWantToTalk}
```
