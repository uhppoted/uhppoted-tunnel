# FAQ

1. How do I reduce the number of log messages?

   Set the log level to _warn_ or _error_:
   - in `console` mode, set the --log-level <debug|info|warn|error> command line argument, e.g.
```
   uhppoted-tunnel --console --log-level warn ...
```

   - in `service` mode, edit the service configuration to set the appropriate level with the `--log-level <debug|info|warn|error>` command line argument

2. In `service` mode how do redirect the log messages to _/dev/null_ ?

   Really not recommend because if anything goes wrong there will nothing whatsoever to work 
   with. But if you're dead set on making things difficult - edit the service command line to
   include the `--console` command line argument. All messages will then go to _stdout_ and 
   you can redirect them using the standard _redirect_ or _pipe_ operators.