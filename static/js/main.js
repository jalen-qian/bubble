let vm = new Vue({
    el: '#app',
    data: {
        input: "",
        visible: false,
        todoList: [],
    },
    created: function () {
        this.refreshTodoList();
    },
    methods: {
        //添加待办事项
        addTodo: function () {
            let _this = this;
            if (_this.input == "") {
                _this.$message.error('请输入待办事项');
                return
            }
            _this.sendRequest("POST", "/v1/todo",
                function (data) {
                    if (data.code === 1000) {
                        _this.$message({
                            message: data.message,
                            type: 'success'
                        })
                    } else {
                        _this.$message.error(data.message);
                    }
                    _this.refreshTodoList();
                },
                function (message) {
                    _this.$message.error(message);
                }, JSON.stringify({title: _this.input}));
            _this.input = "";
        },
        //设置某一个待办事项的状态为“已完成”
        completed: function (todo) {
            let _this = this;
            todo.status = true;
            this.sendRequest('PUT', '/v1/todo/' + todo.id, function (data) {
                if (data.code === 1000) {
                    _this.$message({
                        message: data.message,
                        type: 'success'
                    });
                } else {
                    _this.$message.error(data.message);
                }
                _this.refreshTodoList()
            }, function (message) {
                _this.$message.error(message);
                _this.refreshTodoList()
            });
        },
        //设置某一个待办事项的状态为“未完成”
        unfinished: function (todo) {
            let _this = this;
            todo.status = false;
            this.sendRequest('PUT', '/v1/todo/' + todo.id,
                function (data) {
                    if (data.code === 1000) {
                        _this.$message({
                            message: data.message,
                            type: 'success'
                        });
                    } else {
                        _this.$message.error(data.message);
                    }
                    _this.refreshTodoList()
                }, function (errMsg) {
                    _this.$message.error(errMsg);
                    _this.refreshTodoList()
                })
        },
        deleteTodo: function (todo) {
            let _this = this;
            _this.sendRequest('DELETE', '/v1/todo/' + todo.id,
                function (data) {
                    if (data.code === 1000) {
                        _this.$message({
                            message: data.message,
                            type: 'success'
                        });
                        _this.refreshTodoList()
                    } else {
                        _this.$message.error(data.message);
                    }
                })
        },
        //刷新当前的待办事项列表
        refreshTodoList: function () {
            let _this = this;
            this.sendRequest("GET", "/v1/todo", function (data) {
                if (data.code === 1000) {
                    _this.todoList = data.data;
                } else {
                    _this.$message({
                        message: data.message,
                        type: 'danger'
                    });
                }
            })
        },
        sendRequest: function (method, url, success, error, data) {
            //第一步：建立所需的对象
            let httpRequest = new XMLHttpRequest();
            //第二步：打开连接  将请求参数写在url中
            httpRequest.open(method, url, true);
            if (method == "POST") {
                httpRequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            }
            httpRequest.send(data ? data : null);//第三步：发送请求  将请求参数写在URL中
            /**
             * 获取数据后的处理程序
             */
            httpRequest.onreadystatechange = function () {
                if (httpRequest.readyState === 4 && httpRequest.status === 200) {
                    let data = httpRequest.responseText;//获取到json字符串，还需解析
                    try {
                        data = JSON.parse(data);
                        if (typeof success == "function") {
                            success(data)
                        }
                    } catch (e) {
                        if (typeof error == "function") {
                            error("响应数据异常")
                        }
                    }

                } else if (httpRequest.status !== 200) {
                    if (typeof error == "function") {
                        error("请求失败，请检查网络是否正常！")
                    }
                }
            };
        },
    }
});