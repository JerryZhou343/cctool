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
  final TextStyle headStyle = Theme.of(context).textTheme.headline4;
  //var titleStyle = TextStyle(fontSize: 20, decoration: TextDecoration.none);
  return SimpleDialog(
    title:
        Title(color: Colors.lightBlue, child: Text("创建翻译任务", style: headStyle)),
    children: <Widget>[
      Form(
          child: Column(children: <Widget>[
        TextFormField(
          initialValue: task.from,
          decoration: InputDecoration(
            prefixText: '源语言：',
            prefixStyle: style,
          ),
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
          initialValue: task.to,
          decoration: InputDecoration(prefixStyle: style, prefixText: "目标语言："),
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
        CheckboxListTile(
          title: const Text('双语'),
          activeColor: Colors.blue,
          value: true,
          onChanged: (bool value) {
            task.merge = value;
          },
        ),
      ])),
    ],
  );
}
