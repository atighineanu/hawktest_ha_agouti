# hawktest\_ha\_agouti
####What it does:
	* logs into the ha_cluster
	* sets stonith-sbd to maintenance mode 
	* checks from crm if it is so
	* sets back stonith-sbd to normal mode
	* clears the state of stonith-sbd
	* clears the state of 1st listed node
	* sets the 1st node to maintenance mode 
	* checks from crm if it is so, after sets back
	* submits a hawk history log
	* creates a new primitive (cool_primitive), with 
	  following properties:
	  			- class: ocf
	  			- type: anything
	  			- binfile: file
	  			- start: (timeout) 35s
	  			- stop: (timeout) 15s, stop on-fail
	  			- monitoring: (timeout) 9s, (interval) 13s
	* checks if the primitive cool_primitive is indeed set with right parameters (from 'crm  config show')
	* deletes the primitive cool_primitive
	* checks if the primitive is properly deleted (from 'crm resources list')
	END.


####Additional remarks & reqs:
- 1) The test is SLE 15 SP0 (possibly SP1) compatible (so far).
- 2) For the test to run with full coverage, one needs to install the ssh keys from the client machine (the machine where you will run the test from) on the server you're going to test. 
	
	

