## linux security

Linux has several forms of managing security.

One is obviously user id and group id based where the resources are protected from access based on uid and gid.

### SELinux

In RHEL based system one way to protect resources is via selinux.
If selinux is enabled on a system, all system calls will be processed through selinux.
And based on the policies either the call will be allowed or not.

```
sestatus                      # This checks whether selinux is emnabled or not
setenforce 0                  # Put selinux into permissive mode ie only log things but not prevent access
getenforce                    # Get the selinux mode
setenforce 1                  # Re-enable selinux mode back on
cat /etc/sysconfig/selinux    # Contains selinux settings - set the mode here for persisting across reboots
```

selinux policies and settings are only read during boot time from `/etc/sysconfig/selinux`.

selinux uses context to identify the security related information around an object.
To get selinux context use the `Z` flag with commands.
```
ls -lahZ /etc         # To view the security context of files
netstat -tupleZ       # To view the security context of ports
ps auxZ               # To view the security context of running processes
```

A selinux security context looks like this - user_part:role_part:type_part:sensitivity_part
```
system_u:system_r:kernel_t:s0
system_u:system_r:sshd_t:s0
system_u:object_r:admin_home_t:s0
```

selinux user accounts are different from linux user accounts. You cant create them through the cli.
They are defnined in the policy file and only loaded in memory at startup.
These users do not have server logins. They end with `_u`.
Several selinux users can be mapped to a single selinux user.

selinux roles define what a user can do in a certain domain. They are also defined in the policy files.
selinux users can have multiple roles but only one role can be active and that too it could just be limited to a few selinux resource types.

selinux types part defines what kind of file the resource is identified as, ie directory, file, network file etc.

selinux sensitivity part defines the security behavior. sensitivity levels vary between s0 to s3 with s3 being top secret.
The security levels in each sensitivity level can then range from c1 to c1023.

When new objects like files are created selinux automatically assigns security context to it based on location.
When you copy a file over to a new location a new security context is assigned based on the location.
When you move a file the original security context is preserved.

To view selinux policies
```
semanage user -l
```

To change the type context in selinux for a file at runtime
```
chcon -v -t httpd_sys_content_t foo.html         # Change the type context in selinux security context for a file with verbose output
```

Based on the location of a file you can restore the selinux context of a file to what it should be in that location
```
restorecon foo.html                   # Restore context based on current location
restorecon -R -i -v /var/www/html     # If httpd is running from default location, restore the security context of all files recursively for the /var/www/html directory
```

If you want to change the default root of the directory where the httpd program runs from,
other than changing the configuration for httpd, if selinux is enabled you also need to change the selinux policy for the new root.
```
semanage fcontext  -a –t httpd_sys_content_t "/var/tmp/www(/.*)?"
```
The `-a` option is to specify that we want to add a new context.
The fcontext option in the command above tells the main command that we want to change the security context of files

And consecutively, to load the new policies for the files
```
restorecon –R –i /var/tmp/www/
```

selinux stores all contexts in `/etc/selinux/targeted/contexts`. You can view the settings there.


### setuid

Linux can have privileged and unprivileged processes.
Some binaries need root permissions to work, for example to do things like opening sockets etc.

You can notice the setuid bit set in the ping binary for example.
```
$ ls -l /bin/ping
 -rwsr-xr-x 1 root root 44168 May  7 23:51 /bin/ping
```

To remove the setuid bit
```
chmod u-s /bin/ping
```
This removes the setuid bit from the ping binary and the ping binary now wont be able to work for non-root users.

However, for normal user to be able to effectively use the ping command, after we removed the setuid bit
we can enable the `cap_net_raw` linux capability. This will allow the ping binary to open sockets and send ICMP packets,
ie what it needs to perform it's work. But it will not provide root permissions for doing anything else.
So, even if there was a bug in the ping binary it will be difficult to exploit that.

A linux capability can be applied on a binary in one of 3 ways - `effective`, `permitted` or `inheritable`.
These are essentially the 3 different sets of capabilities a binary may have. And various capabilities will be added to one of these 3 sets.
  - effective : the capabilities used by the kernel to perform permission checks for the thread
  - permitted : the capabilities that the thread may assume
  - inheritable : the capabilities inherited by a child proc from the parent proc when forked.
                  If a binary has the `cap_setcap` capability it can change the capabilities of other threads.
                  Using capset command a thread could change it's own capabilities as well.

So, to apply the `cap_net_raw` capability to the list of permitted capabilities of the ping binary
```
setcap cap_net_raw+p /bin/ping
```

To list all capabilities of a binary
```
getcap /bin/ping
```

The execve system call can grant a newly-started program privileges that its parent did not have.
This can be dangerous. To reduce the attack surface, a generic mechanism has been provided - set `no_new_privs` bit.
Once this is set on a proc, it is inherited across fork, clone and execve and this flag cant be unset.
This way forked procs wont be able to alter the setuid bit or add new capabilities.
