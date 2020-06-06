import 'package:cctool/model/translate_task.dart';
import 'package:flutter/material.dart';
import 'package:cctool/widgets/form_check_box.dart';
import 'package:cctool/utils/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

Widget buildTranslateTask(BuildContext context) {
  var task = TranslateTask();
  task.from = "en";
  task.to = "zh";
  task.merge = false;

  String _extension = ".srt,.mp4";
  FileType _pickingType = FileType.custom;

  var style = TextStyle(fontSize: 16, decoration: TextDecoration.none);
  final TextStyle headStyle = Theme.of(context).textTheme.headline4;

  var checkBox = CheckboxFormField(
    title: Title(color: Colors.lightBlue, child: Text("双语", style: style)),
    context: context,
    onSaved: (bool value) {
      task.merge = value;
    },
  );

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
            prefixText: '原始语言：',
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
        checkBox,
        Container(
            padding: EdgeInsets.all(16),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: <Widget>[
                Expanded(
                    child: TextFormField(
                  initialValue: task.files.toString(),
                  decoration:
                      InputDecoration(prefixStyle: style, prefixText: "原始文件:"),
                  validator: (value) {
                    return value.trim().length > 0 ? null : "不能为空";
                  },
                  //当 Form 表单调用保存方法 Save时回调的函数。
                  onSaved: (value) {
                    task.files.add(value);
                  },
                  // 当用户确定已经完成编辑时触发
                  onFieldSubmitted: (value) {},
                )),
                FlatButton(
                    onPressed: () async {
                      try {
                        //key: 是文件名
                        //value: 绝对路径
                        var files = await FilePicker.getMultiFilePath(
                            type: _pickingType,
                            allowedExtensions: (_extension?.isNotEmpty ?? false)
                                ? _extension?.replaceAll(' ', '')?.split(',')
                                : null);
                        files.forEach((key, value) {
                          task.files.add(value);
                        });
                      } catch (e) {
                        print("Unsupported operation" + e.toString());
                      }
                    },
                    child: Text("浏览文件"))
              ],
            )),
        Container(
          padding: EdgeInsets.all(16),
          child: Row(
            children: <Widget>[
              Expanded(
                child: RaisedButton(
                  padding: EdgeInsets.all(15),
                  child: Text(
                    "提交",
                    style: TextStyle(fontSize: 18),
                  ),
                  textColor: Colors.white,
                  color: Theme.of(context).primaryColor,
                  //onPressed: login,
                ),
              ),
            ],
          ),
        )
      ])),
    ],
  );
}
