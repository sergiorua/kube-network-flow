echo 'require k8s.io/kubernetes release-1.16'
echo 'replace ('
curl -sLk https://github.com/kubernetes/kubernetes/raw/release-1.16/go.mod \
| /usr/local/bin/gsed -n -r 's# \./staging/(.*)$# k8s.io/kubernetes/staging/\1 8c3b7d7679ccf368#p'
echo ')'
