package com.github.sumanmukherjee03;

import java.awt.Dimension;
import java.awt.FlowLayout;

import javax.swing.JButton;
import javax.swing.JFrame;
import javax.swing.JLabel;
import javax.swing.JTextField;
import javax.swing.WindowConstants;

public class Main {

  public static void main(String... args) {
    JFrame application = createGUI();
    application.setDefaultCloseOperation(WindowConstants.EXIT_ON_CLOSE);
    application.setVisible(true);
  }

  private static JFrame createGUI() {
    JTextField input = new JTextField();
    input.setPreferredSize(new Dimension(300, 40));

    JButton convertButton = new JButton("Convert");

    JLabel output = new JLabel();
    output.setPreferredSize(new Dimension(300, 40));

    convertButton.addActionListener(event -> {
      output.setText(TitlecaseConverter.convertTextToTitleCase(input.getText()));
    });

    JFrame gui = new JFrame("Titlecase Converter");
    gui.setLayout(new FlowLayout());
    gui.add(input);
    gui.add(convertButton);
    gui.add(output);
    gui.pack();
    gui.setLocationRelativeTo(null); // Centering the screen

    return gui;
  }
}
