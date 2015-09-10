# revel app to do post-training manipulation for word2vec

1. The assumption is you have a binary word2vec training result file.
  * If you don't know word2vec, take a look [here](https://code.google.com/p/word2vec/)
  * You can put the binary word vector file anywhere the revel app has access. /public folder is a good place
2. Change the init.go file in controllers folder according to your path of the word vector file
3. 'revel run wordapp' should spin up the web service and ready to take queries
4. Currently, there are 3 types of queries.
  1. single word query
  2. phrase query where we use underscore to concatenate words, e.g., data\_mining. 
  3. compound query like Google,China;USA. Positive words are the ones before ';', and negative words are after ';' (words are seperated by ','). It gives similar words to sum(positive vectors) - sum(negative vectors). In this example, it means to find a company in China which is like the Google in the USA. The company Baidu or Tencent should be on top if your corpus is large and good. 
