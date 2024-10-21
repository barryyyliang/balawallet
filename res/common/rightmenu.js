
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
        var rec = document.getElementById("tools").getBoundingClientRect()
        if(e.clientX>rec.left && e.clientX <rec.right && e.clientY >rec.top && e.clientY < rec.bottom) {

        } else {          
            document.getElementById("toolMenu").style.height=0;
        }
        var rec2=document.getElementById("system").getBoundingClientRect()
        if(e.clientX>rec2.left && e.clientX <rec2.right && e.clientY >rec2.top && e.clientY < rec2.bottom) {

        } else {
            document.getElementById("systemMenu").style.height=0;
        }        
    }

    function onMenuOver(id){
        if (id != "tools"){
            document.getElementById("toolMenu").style.height=0;
        }
        if (id != "system"){
            document.getElementById("systemMenu").style.height=0;
        }  
    }

    function onSystemOut(){
        var e = event || window.event
        var rec = document.getElementById("systemMenu").getBoundingClientRect()
        if(e.clientX>rec.left && e.clientX <rec.right && e.clientY >rec.top && e.clientY < rec.bottom) {

        } else {          
            document.getElementById("systemMenu").style.height=0;
        }
    }


    function onToolsMenu(){
        var menu=document.getElementById("toolMenu")
         //根据事件对象中鼠标点击的位置，进行定位
        var a = document.getElementById("tools")
        rec = a.getBoundingClientRect()
        menu.style.left=rec.left+'px';
        menu.style.top=rec.bottom+'px';
        //改变自定义菜单的宽，让它显示出来
        menu.style.width='auto';
        menu.style.height='auto';
    }

    function onSystemMenu(){
        var menu=document.getElementById("systemMenu")
         //根据事件对象中鼠标点击的位置，进行定位
        var a = document.getElementById("system")
        rec = a.getBoundingClientRect()
        menu.style.left=rec.left+'px';
        menu.style.top=rec.bottom+'px';
        //改变自定义菜单的宽，让它显示出来
        menu.style.width='auto';
        menu.style.height='auto';
    }

    //支持 ctrl c  和 ctrl v 快捷键
    document.onkeydown = function(event) {
        if(event.ctrlKey || event.metaKey){            
            if (event.key === 'v' || event.key === 'V'){
                MYMENU.paste()
            } else if (event.key === 'c' || event.key === 'C') {
                MYMENU.copy()
            }
        } 
    }
 
    var MYMENU = {
        copy:function (){
            var userSelection= window.getSelection().toString()
            if(activeArea != undefined){
                if (activeArea.selectionStart != undefined) {
                     userSelection = activeArea.value.substring(
                        activeArea.selectionStart, activeArea.selectionEnd
                    )
                }
            }
  
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
            if (navigator.clipboard) {
                // clipboard api 复制
                navigator.clipboard.readText().then( (result) => {
                    if(activeArea != undefined){
                        activeArea.nodeValue = result
                    }
                }
                )
            } else {
                if(activeArea != undefined){
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
            }
        }
    }

    var MYDIALOG = {
        //模态对话框
        ShowDiv:function (dialog_content){
            var showdiv = document.getElementById("MyDiv");
            showdiv.style.display='block';
            var bgdiv = document.getElementById("fade");
            bgdiv.style.display='block';
            bgdiv.style.width = document.body.scrollWidth;
            bgdiv.style.height = document.body.scrollHeight;
            var contentdiv = document.getElementById("dialogcontent");
            contentdiv.innerHTML = dialog_content
            //$("#"+"fade").height($(document).height());
        },

        //关闭模态对话框
        CloseDiv:function (){
            document.getElementById("MyDiv").style.display='none';
            document.getElementById("fade").style.display='none';
        }
    }

