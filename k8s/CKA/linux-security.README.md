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

selinux roles define what a user can do in a certain domain. They are also defined in the policy files.

selinux types part defines what kind of file the resource is identified as, ie directory, file, network file etc.

selinux sensitivity part defines the security behavior. security levels vary between c0 to c3 with c3 being top secret.
The levels in c3 range can go up to c1023.

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
The fcontext option tells the main command that we want to change the security context of files

And consecutively, to load the new policies for the files
```
restorecon –R –i /var/tmp/www/
```

selinux stores all contexts in `/etc/selinux/targeted/contexts`. You can view the settings there.
