# kube-apparmor-manager
Manage AppAmormor profiles for Kubernetes cluster

## Behind the Scenes
- `AppArmorProfile` CRD is created and `AppArmorProfile` objects are stored in etcd.
- Actual AppArmor profiles will be created(updated) acorss all worker nodes through synchronizing with `AppArmorProfile` objects.

### AppArmorProfile Object Explained
```
apiVersion: crd.security.sysdig.com/v1alpha1
kind: AppArmorProfile
metadata:
  name: apparmorprofile-sample
spec:
  rules: |
    # This is the default deny mode of AppArmor profile.
    # List the allow rules here separated by new line character.
    
    # allow few read/write activities
    allow /etc/* r,
    allow /tmp/* rw,

    # allow few commands execution
    allow /bin/echo mrix,
    allow /bin/sleep mrix,
    allow /bin/cat mrix,
  enforced: true # set profile to enforcement mode if true (complain mode if false)
```

## Configure Environment
- `SSH_USERNAME`: SSH username to access worker nodes (default: admin)
- `SSH_PERM_FILE`: SSH private key to access worker ndoes (default: $HOME/.ssh/id_rsa)
- `SSH_PASSPHRASE`: SSH passphrase (only applicable if the private key is passphrase protected)

## Usage
```
Usage:
  kube-apparmor-manager [command]

Available Commands:
  enabled     Check AppArmor status on worker nodes
  enforced    Check AppArmor profile enforcement status on worker nodes
  help        Help about any command
  init        Install CRD in the cluster and AppArmor services on worker nodes
  sync        Synchronize the AppArmor profiles from the Kubernetes database (etcd) to worker nodes
```

## Example Output

### AppArmor enabled status
```
$ ./kube-apparmor-manager enabled
+-------------------------------+---------------+----------------+--------+------------------+
|           NODE NAME           |  INTERNAL IP  |  EXTERNAL IP   |  ROLE  | APPARMOR ENABLED |
+-------------------------------+---------------+----------------+--------+------------------+
| ip-172-20-45-132.ec2.internal | 172.20.45.132 | 54.91.xxx.xx   | master | false            |
| ip-172-20-54-2.ec2.internal   | 172.20.54.2   | 54.82.xx.xx    | node   | true             |
| ip-172-20-58-7.ec2.internal   | 172.20.58.7   | 18.212.xxx.xxx | node   | true             |
+-------------------------------+---------------+----------------+--------+------------------+
```

### AppArmor enforced profiles
```
./kube-apparmor-manager enforced
+-------------------------------+--------+------------------------------------------------------+
|           NODE NAME           |  ROLE  |                  ENFORCED PROFILES                   |
+-------------------------------+--------+------------------------------------------------------+
| ip-172-20-45-132.ec2.internal | master |                                                      |
| ip-172-20-54-2.ec2.internal   | node   | /usr/sbin/ntpd,apparmorprofile-sample,docker-default |
| ip-172-20-58-7.ec2.internal   | node   | /usr/sbin/ntpd,apparmorprofile-sample,docker-default |
+-------------------------------+--------+------------------------------------------------------+
```

### Sync

When ever there is change to `AppArmorProfile` object, run `sync` to synchronize across all the worker nodes.
```
$ ./kube-apparmor-manager sync
**** Host: 54.82.xx.xx:22 ****
** Execute command: echo 'profile apparmorprofile-sample flags=(attach_disconnected) {
	allow /etc/* r,
	allow /tmp/* rw,
	allow /bin/echo mrix,
	allow /bin/sleep mrix,
	allow /bin/cat mrix,
}' > /tmp/apparmorprofile-sample **

** Execute command: mv /tmp/apparmorprofile-sample /etc/apparmor.d/apparmorprofile-sample **

**** Host: 54.82.xx.xx:22 ****
** Execute command: aa-enforce /etc/apparmor.d/apparmorprofile-sample **
Setting /etc/apparmor.d/apparmorprofile-sample to enforce mode.

**** Host: 18.212.xxx.xxx:22 ****
** Execute command: echo 'profile apparmorprofile-sample flags=(attach_disconnected) {
	allow /etc/* r,
	allow /tmp/* rw,
	allow /bin/echo mrix,
	allow /bin/sleep mrix,
	allow /bin/cat mrix,
}' > /tmp/apparmorprofile-sample **

** Execute command: mv /tmp/apparmorprofile-sample /etc/apparmor.d/apparmorprofile-sample **

**** Host: 18.212.xxx.xxx:22 ****
** Execute command: aa-enforce /etc/apparmor.d/apparmorprofile-sample **
Setting /etc/apparmor.d/apparmorprofile-sample to enforce mode.
```
