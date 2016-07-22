# java-shell-deploy

Toolbox for continuous deployment of a Java application.

The dispatcher is a simple golang app that listens to an url (for example a Github Webhook https://developer.github.com/webhooks/) to trigger an action.

The service shell script manages the start/stop/build actions and is based on Gustavo Straube's approach (http://gustavostraube.wordpress.com/2009/11/05/writing-an-init-script-for-a-java-application/).
