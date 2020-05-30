import 'package:cctool/model/translate_task.dart';
import 'package:flutter/material.dart';
import 'package:cctool/model/choice.dart';
import 'package:cctool/common/const.dart';

Widget buildTranslateTask(BuildContext context) {
  var task = TranslateTask();
  task.from = "en";
  task.to = "zh";
  task.merge = false;

  var style = TextStyle(fontSize: 16, decoration: TextDecoration.none);
  return Center(
    child:Form(child: Column(
      children: <Widget>[
        TextFormField(
          decoration: InputDecoration(
            labelText: '源语言',
            //hintText: "用户名或邮箱",
            hintStyle: TextStyle(
              color: Colors.grey,
              fontSize: 13,
            ),
            //prefixIcon: Icon(Icons.person),
          ),
          //校验用户
          validator: (value) {
            return value.trim().length > 0 ? null : "不能为空";
          },
          //当 Form 表单调用保存方法 Save时回调的函数。
          onSaved: (value) {
            task.from = value;
          },
          // 当用户确定已经完成编辑时触发
          onFieldSubmitted: (value) {},
        ),
        TextFormField(
          decoration: InputDecoration(
            labelText: '目标语言',
            //hintText: '你的登录密码',
            hintStyle: TextStyle(
              color: Colors.grey,
              fontSize: 13,
            ),
            //prefixIcon: Icon(Icons.lock),
          ),
          //是否是密码
          obscureText: false,
          //校验密码
          validator: (value) {
            return value.trim().length > 0 ? null : "不能为空";
          },
          onSaved: (value) {
            task.to = value;
          },
        )
      ],
    )));
}
