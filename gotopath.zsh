gotopath_() {
	# Serve mode
	if [ "$1" = "-serve" ]; then
	 	gotopath -serve
		return $?
	else
		gopath=`gotopath $@`
		ret=$?
		if [ $ret -eq 0 ]; then
			echo "Go to $gopath"
			cd $gopath
		else
			echo "Can't go to '$1'"
		fi
	fi
}

alias goto='gotopath_'
alias g='gotopath_'
