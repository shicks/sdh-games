# sdh-games

Basic framework for mostly peer-to-peer web games.
The general approach is to have a very simple
AppEngine server to mediate p2p communication
between players.  Effectively this provides basic
chat as well as fair card shuffling and dice
rolling.  All clients will be able to ask for
any shuffled card, but requests will be broadcast
so that cheating is immediately detectable.

We will also provide a framework for writing bots.
