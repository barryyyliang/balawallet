
    var activeArea 
    var userSelection
    window.oncontextmenu=function(e){
        //取消默认的浏览器自带右键 很重要！！
        activeArea = document.activeElement
        userSelection = window.getSelection().toString()
        e.preventDefault();
        //获取我们自定义的右键菜单
        var menu=document.querySelector("#menu");
        //根据事件对象中鼠标点击的位置，进行定位
        menu.style.left=(e.clientX+2+window.scrollX)+'px';
        menu.style.top=(e.clientY+2+window.scrollY)+'px';
        //改变自定义菜单的宽，让它显示出来
        menu.style.width='100px';
        menu.style.height='auto';
    }
    //支持 ctrl c  和 ctrl v 快捷键
    document.onkeydown = function(event) {
        if(event.ctrlKey || event.metaKey){            
            if (event.key === 'v' || event.key === 'V'){
                MYMENU.paste()
            } else if (event.key === 'c' || event.key === 'C') {
                userSelection = window.getSelection().toString()
                MYMENU.copy()
            }
        } 
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

    var MYMENU = {
        copy:function (){
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
                document.execCommand('copy');
            }
            activeArea = null
            userSelection = null
        },

        paste:function(){
            if (navigator.clipboard) {
                // clipboard api 复制
                navigator.clipboard.readText().then( (result) => {
                    if(activeArea != undefined){
                        try{
                            activeArea.nodeValue = result
                        } catch(err){

                        }                        
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
            activeArea = null
            userSelection = null
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

        ShowConfirmDiv: async function (dialog_content){
            var showdiv = document.getElementById("MyDiv");
            showdiv.style.display='flex';
            var bgdiv = document.getElementById("fade");
            bgdiv.style.display='block';
            bgdiv.style.width = document.body.scrollWidth;
            bgdiv.style.height = document.body.scrollHeight;
            var contentdiv = document.getElementById("dialogcontent");
            contentdiv.innerHTML = ""

            var div1 = document.createElement("div")
            div1.setAttribute("style","font-size: 16px;display: flex;flex-grow: 1;align-items:center ;justify-content:center;")
            var p = document.createElement("p")
            p.innerText = dialog_content
            div1.append(p)
            contentdiv.append(div1)

            var div2 = document.createElement("div")
            div2.setAttribute("style","width: 100%;text-align:center;margin-bottom: 10px;")            
            var button1 = document.createElement("input")
            button1.setAttribute("type","button")
            button1.setAttribute("value",await gGetMainRes("main_ok"))
            button1.setAttribute("class","button1")
            button1.setAttribute("style","width:100px;")
            div2.append(button1)

            var button2 = document.createElement("input")
            button2.setAttribute("type","button")
            button2.setAttribute("value",await gGetMainRes("main_cancel"))
            button2.setAttribute("class","button1")
            button2.setAttribute("style","width:100px")
            div2.append(button2)
            contentdiv.append(div2)


            return new Promise((resolve, reject) => {
                button1.onclick = () => {
                    MYDIALOG.CloseDiv();
                    resolve(true);
                };

                button2.onclick = () => {
                    MYDIALOG.CloseDiv();
                    resolve(false)
                }
            });

            //$("#"+"fade").height($(document).height());
        },

        //关闭模态对话框
        CloseDiv:function (){
            document.getElementById("MyDiv").style.display='none';
            document.getElementById("fade").style.display='none';
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
        var e = Event || window.Event
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
        menu.style.left=(rec.left+window.scrollX)+'px';
        menu.style.top=(rec.bottom+window.scrollY)+'px';
        //改变自定义菜单的宽，让它显示出来
        menu.style.width='auto';
        menu.style.height='auto';
    }

    function onSystemMenu(){
        var menu=document.getElementById("systemMenu")
         //根据事件对象中鼠标点击的位置，进行定位
        var a = document.getElementById("system")
        rec = a.getBoundingClientRect()
        menu.style.left=(rec.left+window.scrollX)+'px';
        menu.style.top=(rec.bottom+window.scrollY)+'px';
        //改变自定义菜单的宽，让它显示出来
        menu.style.width='auto';
        menu.style.height='auto';
    }

