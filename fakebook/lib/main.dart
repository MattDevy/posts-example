import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

Future<GetPostsResponse> fetchPosts() async {
  final response =
      await http.get(Uri.http('10.0.2.2:3000', 'api/v1/posts'), headers: {
    "Access-Control-Allow-Origin": "*",
  });

  if (response.statusCode == 200) {
    // If the server did return a 200 OK response,
    // then parse the JSON.
    return GetPostsResponse.fromJson(jsonDecode(response.body));
  } else {
    // If the server did not return a 200 OK response,
    // then throw an exception.
    throw Exception('Failed to load post');
  }
}

class GetPostsResponse {
  final List<Post> posts;
  final String next_page_token;

  GetPostsResponse({this.posts, this.next_page_token});

  factory GetPostsResponse.fromJson(Map<String, dynamic> json) {
    var posts = json['posts'] as List;
    List<Post> postsList = posts.map((i) => Post.fromJson(i)).toList();

    return GetPostsResponse(
      posts: postsList,
      next_page_token: json['next_page_token'],
    );
  }
}

class Post {
  final int id;
  final String title;
  final String author;
  final String text;

  Post({this.id, this.title, this.author, this.text});

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      id: json['id'],
      title: json['title'],
      author: json['author'],
      text: json['text'],
    );
  }
}

void main() => runApp(MyApp());

class MyApp extends StatefulWidget {
  MyApp({Key key}) : super(key: key);

  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  Future<GetPostsResponse> futurePosts;
  final _posts = <Post>[];
  final _biggerFont = TextStyle(fontSize: 18.0);

  @override
  void initState() {
    super.initState();
    futurePosts = fetchPosts();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Fetch Data Example',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: Scaffold(
        appBar: AppBar(
          title: Text('Fetch Data Example'),
        ),
        body: Center(
          child: FutureBuilder<GetPostsResponse>(
            future: futurePosts,
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                return ListView.builder(
                    padding: EdgeInsets.all(16.0),
                    itemBuilder: (context, i) {
                      if (i.isOdd) return Divider();

                      final index = i ~/ 2;
                      if (index >= _posts.length) {
                        _posts.addAll(snapshot.data.posts); /*4*/
                      }
                      return _buildrow(_posts[index]);
                    });
              } else if (snapshot.hasError) {
                return Text("${snapshot.error}");
              }

              // By default, show a loading spinner.
              return CircularProgressIndicator();
            },
          ),
        ),
      ),
    );
  }

  Widget _buildrow(Post post) {
    return ListTile(
        title: Text(
      post.title,
      style: _biggerFont,
    ));
  }
}
