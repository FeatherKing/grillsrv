package com.smokerapp;

import android.content.Intent;
import android.content.SharedPreferences;
import android.net.Uri;
import android.os.Bundle;
import android.preference.PreferenceManager;
import android.support.design.widget.FloatingActionButton;
import android.support.design.widget.Snackbar;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class SettingsActivity extends AppCompatActivity {
    Button button;
    EditText ipDnsText;
    EditText portText;
    EditText usernameText;
    EditText passwordText;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main2);
        ipDnsText = (EditText) findViewById(R.id.grill_ip_text_box);
        portText = (EditText) findViewById(R.id.port_text_box);
        usernameText = (EditText) findViewById(R.id.username);
        passwordText = (EditText) findViewById(R.id.password);
        button = (Button) findViewById(R.id.button1);
//        button.setOnClickListener((View.OnClickListener) this);
        addListenerOnButton();
        loadSavedPreferences();
//        addListenerOnButton();

    }

    private void loadSavedPreferences() {
        SharedPreferences sharedPreferences = PreferenceManager
                .getDefaultSharedPreferences(this);
        String ip_dns = sharedPreferences.getString("ip_dns", "");
        ipDnsText.setText(ip_dns);
        String port = sharedPreferences.getString("port", "");
        portText.setText(port);
        String username = sharedPreferences.getString("username", "");
        usernameText.setText(username);
        String password = sharedPreferences.getString("password", "");
        passwordText.setText(password);
    }
    private void savePreferences(String key, String value) {
        SharedPreferences sharedPreferences = PreferenceManager
                .getDefaultSharedPreferences(this);
        SharedPreferences.Editor editor = sharedPreferences.edit();
        editor.putString(key, value);
        editor.commit();
    }


    public void addListenerOnButton() {

        button = (Button) findViewById(R.id.button1);

        button.setOnClickListener(new View.OnClickListener() {

            @Override
            public void onClick(View view) {
                savePreferences("ip_dns", ipDnsText.getText().toString());
                savePreferences("port", portText.getText().toString());
                savePreferences("username", usernameText.getText().toString());
                savePreferences("password", passwordText.getText().toString());
                Snackbar.make(view, "Settings Saved", Snackbar.LENGTH_LONG)
                        .setAction("Action", null).show();
//                finish();



            }

        });

    }

}
