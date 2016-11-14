# Gotopath client
gotopath()
{
	gopath=`command gotopath -request $@`
	ret=$?
	if [ $ret -eq 0 ]; then
		echo "Go to $gopath"
		cd $gopath
	else
		echo "Can't go to '$1'"
	fi
}

alias goto='gotopath'
alias g='gotopath'
