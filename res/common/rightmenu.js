
    var activeArea 
    window.oncontextmenu=function(e){
        //取消默认的浏览器自带右键 很重要！！
        activeArea = document.activeElement
        e.preventDefault();
        //获取我们自定义的右键菜单
        var menu=document.querySelector("#menu");
        //根据事件对象中鼠标点击的位置，进行定位
        menu.style.left=e.clientX+'px';
        menu.style.top=e.clientY+'px';
        //改变自定义菜单的宽，让它显示出来
        menu.style.width='100px';
        menu.style.height='auto';
    }
    //鼠标左键任意位置单击, 关闭右键菜单
    window.onclick=function(e){
        //用户触发click事件就可以关闭了，因为绑定在window上，按事件冒泡处理，不会影响菜单的功能
        document.querySelector('#menu').style.height=0;
    }
 
    var MYMENU = {
        copy:function (){
            var userSelection= window.getSelection().toString()
            window.doWebLog("select:"+ userSelection)
            window.doWebLog("nodename:"+ window.getSelection())
            window.doWebLog("选择copy")
            if(activeArea != undefined){
                window.doWebLog("选中内容:"+activeArea.nodeName)
                if (activeArea.selectionStart != undefined) {
                    window.doWebLog("选中部分内容")
                     userSelection = activeArea.value.substring(
                        activeArea.selectionStart, activeArea.selectionEnd
                    )
                }
            }

            window.doWebLog("content:"+userSelection)
    
            if (navigator.clipboard) {
                // clipboard api 复制
                navigator.clipboard.writeText(userSelection);
            } else {
                var textarea = document.getElementById('copyarea');
                // 赋值
                textarea.value = userSelection;
                // 选中
                textarea.select();
                // 复制
                document.execCommand('copy', true);
            }
        },

        paste:function(){
            window.doWebLog("paste")
            if (navigator.clipboard) {
                // clipboard api 复制
                window.doWebLog("get clip")
                navigator.clipboard.readText().then( (result) => {
                    window.doWebLog("get paste")
                    if(activeArea != undefined){
                        window.doWebLog(activeArea.nodeName)
                        activeArea.nodeValue = result
                    }
                }
                )
            } else {
                window.doWebLog("not clip")
                if(activeArea != undefined){
                    window.doWebLog(activeArea.nodeName)
                    window.doWebLog(activeArea.nodeValue)
                    window.doWebLog(activeArea.id)

                    activeArea.focus()
                }
                //text = window.Clipboard.getData('Text');
                //var textarea = document.getElementById('copyarea');
                // 赋值
                //textarea.value = "";
                // 选中
                //textarea.focus();
                // 复制
                document.execCommand('paste');

                window.doWebLog(document.getElementById('copyarea').value)
            }
        }
    }

    var MYDIALOG = {
        //模态对话框
        ShowDiv:function (show_div, bg_div, dialog_content){
            var showdiv = document.getElementById(show_div);
            showdiv.style.display='block';
            var bgdiv = document.getElementById(bg_div);
            bgdiv.style.display='block';
            bgdiv.style.width = document.body.scrollWidth;
            bgdiv.style.height = document.body.scrollHeight;
            var contentdiv = document.getElementById("dialogcontent");
            contentdiv.innerHTML = dialog_content
            //$("#"+bg_div).height($(document).height());
        },

        //关闭模态对话框
        CloseDiv:function (show_div,bg_div){
            document.getElementById(show_div).style.display='none';
            document.getElementById(bg_div).style.display='none';
        }
    }

