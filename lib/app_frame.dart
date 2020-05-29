import 'package:cctool/widgets/factory.dart';
import 'package:flutter/material.dart';
import 'package:cctool/common/const.dart';
import 'package:cctool/model/choice.dart';
import 'package:cctool/widgets/choice_card.dart';

class AppFrame extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CCTool',
      theme: ThemeData(
        // If the host is missing some fonts, it can cause the
        // text to not be rendered or worse the app might crash.
        fontFamily: 'Roboto',
        primarySwatch: Colors.blue,
      ),
      home: DefaultTabController(
        length: choices.length,
        child: Scaffold(
          appBar: AppBar(
            title: const Text('CCTool'),
            bottom: TabBar(
              isScrollable: true,
              tabs: choices.map<Widget>((Choice choice) {
                return Tab(
                  text: choice.title,
                  icon: Icon(choice.icon),
                );
              }).toList(),
            ),
          ),
          body: TabBarView(
            children: choices.map<Widget>((Choice choice) {
              return Padding(
                padding: const EdgeInsets.all(16.0),
                child: ChoiceCard(choice: choice),
              );
            }).toList(),
          ),
        ),
      ),
    );
  }
}

const List<Choice> choices = <Choice>[
  Choice(title: Translate, icon: Icons.g_translate),
  Choice(title: Merge, icon: Icons.call_merge),
  Choice(title: Generate, icon: Icons.android),
  Choice(title: Clear, icon: Icons.clear_all),
  Choice(title: Convert, icon: Icons.autorenew),
  Choice(title: Setting, icon: Icons.settings_applications),
];
